/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package cli

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terraform-docs/terraform-docs/internal/testutil"
)

func TestFileWriter(t *testing.T) {
	content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
	tests := map[string]struct {
		file     string
		mode     string
		template string
		begin    string
		end      string
		writer   io.Writer

		expected string
		wantErr  bool
		errMsg   string
	}{
		// Successful writes
		"ModeInject": {
			file:     "mode-inject.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   &bytes.Buffer{},

			expected: "mode-inject",
			wantErr:  false,
			errMsg:   "",
		},
		"ModeReplaceWithComment": {
			file:     "mode-replace.md",
			mode:     "replace",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   &bytes.Buffer{},

			expected: "mode-replace-with-comment",
			wantErr:  false,
			errMsg:   "",
		},
		"ModeReplaceWithoutComment": {
			file:     "mode-replace.md",
			mode:     "replace",
			template: outputContent,
			begin:    "",
			end:      "",
			writer:   &bytes.Buffer{},

			expected: "mode-replace-without-comment",
			wantErr:  false,
			errMsg:   "",
		},
		"ModeReplaceWithoutTemplate": {
			file:     "mode-replace.md",
			mode:     "replace",
			template: "",
			begin:    "",
			end:      "",
			writer:   &bytes.Buffer{},

			expected: "mode-replace-without-template",
			wantErr:  false,
			errMsg:   "",
		},

		// Error writes
		"ModeInjectNoFile": {
			file:     "file-missing.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "open testdata/writer/file-missing.md: no such file or directory",
		},
		"EmptyTemplate": {
			file:     "not-applicable.md",
			mode:     "inject",
			template: "",
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "template is missing",
		},
		"EmptyFile": {
			file:     "empty-file.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "file content is empty",
		},
		"BeginCommentMissing": {
			file:     "begin-comment-missing.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "begin comment is missing",
		},
		"EndCommentMissing": {
			file:     "end-comment-missing.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "end comment is missing",
		},
		"EndCommentBeforeBegin": {
			file:     "end-comment-before-begin.md",
			mode:     "inject",
			template: OutputTemplate,
			begin:    outputBeginComment,
			end:      outputEndComment,
			writer:   nil,

			expected: "",
			wantErr:  true,
			errMsg:   "end comment is before begin comment",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			writer := &fileWriter{
				file: tt.file,
				dir:  filepath.Join("testdata", "writer"),

				mode: tt.mode,

				template: tt.template,
				begin:    tt.begin,
				end:      tt.end,

				writer: tt.writer,
			}

			_, err := io.WriteString(writer, content)

			if tt.wantErr {
				assert.NotNil(err)
				assert.Equal(tt.errMsg, err.Error())
			} else {
				assert.Nil(err)

				w, ok := tt.writer.(*bytes.Buffer)
				assert.True(ok, "tt.writer is not a valid bytes.Buffer")

				expected, err := testutil.GetExpected("writer", tt.expected)
				assert.Nil(err)

				assert.Equal(expected, w.String())
			}
		})
	}
}
