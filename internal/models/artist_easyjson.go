// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson75485a89Decode20212LostPointerInternalModels(in *jlexer.Lexer, out *Artists) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Artists, 0, 0)
			} else {
				*out = Artists{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Artist
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson75485a89Encode20212LostPointerInternalModels(out *jwriter.Writer, in Artists) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Artists) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson75485a89Encode20212LostPointerInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Artists) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson75485a89Encode20212LostPointerInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Artists) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson75485a89Decode20212LostPointerInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Artists) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson75485a89Decode20212LostPointerInternalModels(l, v)
}
func easyjson75485a89Decode20212LostPointerInternalModels1(in *jlexer.Lexer, out *Artist) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "video":
			out.Video = string(in.String())
		case "tracks":
			if in.IsNull() {
				in.Skip()
				out.Tracks = nil
			} else {
				in.Delim('[')
				if out.Tracks == nil {
					if !in.IsDelim(']') {
						out.Tracks = make([]Track, 0, 0)
					} else {
						out.Tracks = []Track{}
					}
				} else {
					out.Tracks = (out.Tracks)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Track
					(v4).UnmarshalEasyJSON(in)
					out.Tracks = append(out.Tracks, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "albums":
			if in.IsNull() {
				in.Skip()
				out.Albums = nil
			} else {
				in.Delim('[')
				if out.Albums == nil {
					if !in.IsDelim(']') {
						out.Albums = make([]Album, 0, 0)
					} else {
						out.Albums = []Album{}
					}
				} else {
					out.Albums = (out.Albums)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Album
					(v5).UnmarshalEasyJSON(in)
					out.Albums = append(out.Albums, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson75485a89Encode20212LostPointerInternalModels1(out *jwriter.Writer, in Artist) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Avatar != "" {
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	if in.Video != "" {
		const prefix string = ",\"video\":"
		out.RawString(prefix)
		out.String(string(in.Video))
	}
	if len(in.Tracks) != 0 {
		const prefix string = ",\"tracks\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v6, v7 := range in.Tracks {
				if v6 > 0 {
					out.RawByte(',')
				}
				(v7).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.Albums) != 0 {
		const prefix string = ",\"albums\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v8, v9 := range in.Albums {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Artist) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson75485a89Encode20212LostPointerInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Artist) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson75485a89Encode20212LostPointerInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Artist) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson75485a89Decode20212LostPointerInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Artist) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson75485a89Decode20212LostPointerInternalModels1(l, v)
}