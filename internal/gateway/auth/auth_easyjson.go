// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package handlers

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

func easyjson4a0f95aaDecodeKudagoInternalGatewayAuth(in *jlexer.Lexer, out *UserResponse) {
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
			out.ID = int(in.Int())
		case "username":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "image":
			out.ImageURL = string(in.String())
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
func easyjson4a0f95aaEncodeKudagoInternalGatewayAuth(out *jwriter.Writer, in UserResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"image\":"
		out.RawString(prefix)
		out.String(string(in.ImageURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth(l, v)
}
func easyjson4a0f95aaDecodeKudagoInternalGatewayAuth1(in *jlexer.Lexer, out *RegisterRequest) {
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
		case "username":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson4a0f95aaEncodeKudagoInternalGatewayAuth1(out *jwriter.Writer, in RegisterRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegisterRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegisterRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegisterRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegisterRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth1(l, v)
}
func easyjson4a0f95aaDecodeKudagoInternalGatewayAuth2(in *jlexer.Lexer, out *LoginRequest) {
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
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson4a0f95aaEncodeKudagoInternalGatewayAuth2(out *jwriter.Writer, in LoginRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LoginRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LoginRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LoginRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LoginRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth2(l, v)
}
func easyjson4a0f95aaDecodeKudagoInternalGatewayAuth3(in *jlexer.Lexer, out *AuthResponse) {
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
		case "user":
			(out.User).UnmarshalEasyJSON(in)
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
func easyjson4a0f95aaEncodeKudagoInternalGatewayAuth3(out *jwriter.Writer, in AuthResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		(in.User).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AuthResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AuthResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeKudagoInternalGatewayAuth3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AuthResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AuthResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeKudagoInternalGatewayAuth3(l, v)
}
