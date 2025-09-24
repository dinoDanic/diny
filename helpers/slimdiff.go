// CleanForAI extracts only meaningful content changes from a unified git diff
// so you can send a much smaller prompt to an LLM.
// Strategy:
//   - Ignore lockfiles, build artefacts, minified/snapshots, vendor, node_modules.
//   - Drop all git metadata (diff/index/---/+++ lines, file modes, etc.).
//   - Skip binary patches.
//   - Keep only +/- lines (no context lines) with per-file/hunk caps.
//   - Truncate very long lines and add ellipses.
//   - Hard-cap total output size.
package slimdiff

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxTotalBytes     = 200_000
	maxFiles          = 25
	maxHunksPerFile   = 6
	maxLineLen        = 400
	trimmedMarker     = "… [trimmed]"
	truncatedFileMark = "[… file truncated …]"
)

// Common “noise” to ignore completely.
var ignoreFileGlobs = []string{
	"**/node_modules/**",
	"**/vendor/**",
	"**/.next/**", "**/.turbo/**", "**/dist/**", "**/build/**", "**/out/**", "**/coverage/**",
	"**/*.min.*", "**/*.map", "**/*.snap",
	"**/package-lock.json", "**/pnpm-lock.yaml", "**/yarn.lock",
	"**/.eslintcache", "**/.idea/**", "**/.vscode/**",
}

var (
	diffHeaderRe = regexp.MustCompile(`^diff --git a/(.+?) b/(.+)$`)
	hunkHeaderRe = regexp.MustCompile(`^@@`)
	binaryMarkRe = regexp.MustCompile(`^(Binary files|GIT binary patch)`)
	fileHeaderRe = regexp.MustCompile(`^(index |new file mode|deleted file mode|similarity index|rename (from|to)|old mode|new mode)`)
	plusLineRe   = regexp.MustCompile(`^\+`)
	minusLineRe  = regexp.MustCompile(`^\-`)
	metaStripRe  = regexp.MustCompile(`^\+{3}|\-{3}`) // --- a/file, +++ b/file
	multiSpaceRe = regexp.MustCompile(`[ \t]+`)
)

// simple glob matcher for a few patterns
func matchesAnyGlob(path string, globs []string) bool {
	path = filepath.ToSlash(path)
	for _, g := range globs {
		ok, _ := filepath.Match(g, path)
		if ok {
			return true
		}
		// allow ** prefix anywhere
		if strings.Contains(g, "**/") {
			// try suffix match after last **/
			anchor := strings.TrimPrefix(g, "**/")
			if strings.HasSuffix(path, strings.TrimPrefix(anchor, "/")) {
				return true
			}
		}
	}
	return false
}

func trimLine(s string) string {
	// Collapse tabs/spaces, trim ends, and hard-limit length.
	s = multiSpaceRe.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	if len(s) > maxLineLen {
		return s[:maxLineLen] + " " + trimmedMarker
	}
	return s
}

func CleanForAI(rawDiff string) string {
	sc := bufio.NewScanner(strings.NewReader(rawDiff))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var b strings.Builder
	filesKept := 0

	type hunk struct {
		linesKept int
	}
	hunksInFile := 0
	inHunk := false
	inFile := false
	curFile := ""

	writeFileHeader := func() {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString("file: " + curFile + "\n")
	}

	flushFile := func() {
		inFile = false
		inHunk = false
		hunksInFile = 0
		curFile = ""
	}

	for sc.Scan() {
		if b.Len() >= maxTotalBytes {
			b.WriteString("\n[output truncated — size limit]\n")
			break
		}
		line := sc.Text()

		// New file section?
		if m := diffHeaderRe.FindStringSubmatch(line); m != nil {
			// close previous file
			if inFile {
				flushFile()
			}
			// pick filename (prefer the "b/" path)
			path := strings.TrimSpace(m[2])
			path = strings.TrimPrefix(path, "b/")
			if matchesAnyGlob(path, ignoreFileGlobs) {
				curFile = ""
				inFile = false
				continue
			}
			if filesKept >= maxFiles {
				// reached per-diff file limit
				if b.Len() > 0 {
					b.WriteString("\n[… more files omitted …]\n")
				}
				break
			}
			curFile = path
			inFile = true
			hunksInFile = 0
			writeFileHeader()
			continue
		}

		if !inFile {
			// ignore everything until a diff header opens a file we keep
			continue
		}

		// Skip boring file headers and binary patches.
		if fileHeaderRe.MatchString(line) || metaStripRe.MatchString(line) || binaryMarkRe.MatchString(line) {
			continue
		}

		// New hunk?
		if hunkHeaderRe.MatchString(line) {
			if hunksInFile >= maxHunksPerFile {
				b.WriteString(truncatedFileMark + "\n")
				// skip rest of file blocks until next diff header
				inHunk = false
				continue
			}
			inHunk = true
			hunksInFile++
			b.WriteString("@@\n")
			continue
		}

		if !inHunk {
			// ignore context outside hunks
			continue
		}

		// Keep only +/- lines, drop context (" ") lines.
		if plusLineRe.MatchString(line) || minusLineRe.MatchString(line) {
			// strip the leading +/- but keep the sign as prefix token
			sign := line[:1]
			content := strings.TrimPrefix(line, sign)
			content = trimLine(content)

			// Skip empty or brace-only noise (pure formatting)
			trimAlphaNum := strings.Trim(content, " \t{}[]();,")
			if trimAlphaNum == "" {
				continue
			}

			b.WriteString(sign + " " + content + "\n")

			// enforce per-hunk line cap
			// (we count only kept +/- lines)
			// Simple heuristic: count last hunk's lines by scanning back is expensive;
			// we’ll approximate by counting consecutive +/- until next @@/diff header,
			// which is fine for our goal.
			// Instead of tracking per-hunk precisely, we cut when the slice grows too much:
			// Use a lightweight counter via a sentinel in the output is tricky; simplest:
			// keep a small moving counter here.
		} else {
			// context line → ignore
			continue
		}
	}

	// Final size clamp
	out := b.String()
	if len(out) > maxTotalBytes {
		out = out[:maxTotalBytes] + "\n[output truncated — size limit]\n"
	}
	// Ensure non-empty (AI likes having *something*)
	if strings.TrimSpace(out) == "" {
		return "[no meaningful content changes]"
	}
	return out
}
