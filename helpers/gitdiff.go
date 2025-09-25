package helpers

import (
	"regexp"
	"strings"
)

func OptimizeGitDiff(rawDiff string) string {
	const (
		contextRadius    = 2   // number of context lines to keep around each change
		maxFileDiffLines = 800 // hard cap per file to avoid excessive diffs
	)

	ignoreFilePatterns := []*regexp.Regexp{
		// lock / generated files
		regexp.MustCompile(`(?i)^.*(?:^|/)(?:package-lock\.json|yarn\.lock|pnpm-lock\.yaml|composer\.lock|Cargo\.lock)$`),
		regexp.MustCompile(`(?i)^.*(?:^|/)(?:dist|build|out|coverage|\.next|\.turbo|target|node_modules)(?:/|$)`),
		regexp.MustCompile(`(?i)^.*(?:^|/)__snapshots__(?:/|$)`),

		// map / minified
		regexp.MustCompile(`(?i)\.map$`),
		regexp.MustCompile(`(?i)\.min\.(js|css)$`),

		// binary / assets
		regexp.MustCompile(`(?i)\.(png|jpe?g|gif|webp|svg|ico|pdf|wasm|ttf|otf|woff2?)$`),
		regexp.MustCompile(`(?i)\.(mp4|mov|avi|mkv|webm|mp3|wav|flac)$`),
		regexp.MustCompile(`(?i)\.(zip|gz|bz2|xz|7z|rar|tar)$`),

		// misc
		regexp.MustCompile(`(?i)\.DS_Store$`),
	}

	lines := strings.Split(rawDiff, "\n")
	var out []string

	// regex helpers
	reDiffHeader := regexp.MustCompile(`^diff --git a/(.+?) b/(.+)$`)
	reBinary := regexp.MustCompile(`^Binary files .* differ$`)
	reIndexLine := regexp.MustCompile(`^index [0-9a-f]+\.\.[0-9a-f]+`)
	reOldNew := regexp.MustCompile(`^(\-\-\- a/|\+\+\+ b/)`)
	reHunk := regexp.MustCompile(`^@@ .* @@`)

	// state for each file block
	type fileBlock struct {
		pathB   string
		header  []string
		hunks   [][]string
		ignored bool
	}

	var cur *fileBlock

	shouldIgnore := func(path string) bool {
		for _, rx := range ignoreFilePatterns {
			if rx.MatchString(path) {
				return true
			}
		}
		return false
	}

	flushFile := func(f *fileBlock) {
		if f == nil || f.ignored {
			return
		}
		var fileOut []string
		fileOut = append(fileOut, f.header...)
		total := 0
		for _, h := range f.hunks {
			if total >= maxFileDiffLines {
				fileOut = append(fileOut, "... (truncated)")
				break
			}
			for _, ln := range h {
				fileOut = append(fileOut, ln)
				total++
				if total >= maxFileDiffLines {
					fileOut = append(fileOut, "... (truncated)")
					break
				}
			}
		}
		out = append(out, fileOut...)
	}

	// compress a hunk: keep +/- lines and up to contextRadius context lines
	compressHunk := func(hunk []string) []string {
		if len(hunk) == 0 {
			return hunk
		}
		header := hunk[0]
		body := hunk[1:]

		keep := make([]bool, len(body))

		isChange := func(s string) bool {
			if len(s) == 0 {
				return false
			}
			// + or - but not +++/---
			if strings.HasPrefix(s, "+++") || strings.HasPrefix(s, "---") {
				return false
			}
			return s[0] == '+' || s[0] == '-'
		}

		changeIdx := []int{}
		for i, ln := range body {
			if isChange(ln) {
				changeIdx = append(changeIdx, i)
			}
		}

		if len(changeIdx) == 0 {
			// no changes, return header + first few lines
			limit := len(body)
			if limit > 10 {
				limit = 10
			}
			res := []string{header}
			res = append(res, body[:limit]...)
			if len(body) > limit {
				res = append(res, "...")
			}
			return res
		}

		// mark keep windows around changes
		for _, c := range changeIdx {
			start := c - contextRadius
			if start < 0 {
				start = 0
			}
			end := c + contextRadius
			if end >= len(body) {
				end = len(body) - 1
			}
			for i := start; i <= end; i++ {
				keep[i] = true
			}
		}

		// rebuild hunk with ellipses for skipped sections
		res := []string{header}
		skipping := false
		for i, ln := range body {
			if keep[i] || isChange(ln) || reOldNew.MatchString(ln) {
				res = append(res, ln)
				skipping = false
			} else {
				if !skipping {
					res = append(res, "...")
					skipping = true
				}
			}
		}
		return res
	}

	// main loop
	for i := 0; i < len(lines); i++ {
		ln := lines[i]

		if reDiffHeader.MatchString(ln) {
			// new file starts — flush the previous
			flushFile(cur)

			m := reDiffHeader.FindStringSubmatch(ln)
			pathB := strings.TrimSpace(m[2])
			if strings.HasPrefix(pathB, "b/") && len(pathB) > 2 {
				pathB = pathB[2:]
			}

			cur = &fileBlock{
				pathB:  pathB,
				header: []string{ln},
			}

			// mark ignore if matched
			cur.ignored = shouldIgnore(pathB)
			continue
		}

		if cur == nil {
			continue
		}
		if cur.ignored {
			continue
		}

		if reBinary.MatchString(ln) {
			// ignore binary diff
			cur.ignored = true
			cur.header = nil
			cur.hunks = nil
			continue
		}

		if reIndexLine.MatchString(ln) || reOldNew.MatchString(ln) {
			cur.header = append(cur.header, ln)
			continue
		}

		if reHunk.MatchString(ln) {
			hunk := []string{ln}
			j := i + 1
			for ; j < len(lines); j++ {
				nxt := lines[j]
				if reHunk.MatchString(nxt) || reDiffHeader.MatchString(nxt) {
					break
				}
				hunk = append(hunk, nxt)
			}
			cur.hunks = append(cur.hunks, compressHunk(hunk))
			i = j - 1
			continue
		}

		if strings.HasPrefix(ln, "new file mode ") ||
			strings.HasPrefix(ln, "deleted file mode ") ||
			strings.HasPrefix(ln, "similarity index ") ||
			strings.HasPrefix(ln, "rename from ") ||
			strings.HasPrefix(ln, "rename to ") ||
			strings.HasPrefix(ln, "old mode ") ||
			strings.HasPrefix(ln, "new mode ") {
			cur.header = append(cur.header, ln)
			continue
		}
	}

	flushFile(cur)

	return strings.TrimSpace(strings.Join(out, "\n"))
}
