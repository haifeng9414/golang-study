package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

type Movie struct {
	//Title, Subtitle string
	//Year            int
	Color           bool
	//Actor           map[string]string
	//Oscars          []string
	//Sequel          *string
}

func main() {
	//strangelove := Movie{
	//	Title:    "Dr. Strangelove",
	//	Subtitle: "How I Learned to Stop Worrying and Love the Bomb", Year: 1964,
	//	Color: false,
	//	Actor: map[string]string{
	//		"Dr. Strangelove":            "Peter Sellers",
	//		"Grp. Capt. Lionel Mandrake": "Peter Sellers",
	//		"Pres. Merkin Muffley":       "Peter Sellers",
	//		"Gen. Buck Turgidson":        "George C. Scott",
	//		"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
	//		`Maj. T.J. "King" Kong`:      "Slim Pickens",
	//	},
	//	Oscars: []string{
	//		"Best Actor (Nomin.)",
	//		"Best Adapted Screenplay (Nomin.)", "Best Director (Nomin.)",
	//		"Best Picture (Nomin.)",
	//	},
	//}
	strangelove := Movie{
		//Title:    "Dr. Strangelove",
		//Subtitle: "How I Learned to Stop Worrying and Love the Bomb", Year: 1964,
		Color: false,
		//Actor: map[string]string{
		//	"Dr. Strangelove":            "Peter Sellers",
		//	"Grp. Capt. Lionel Mandrake": "Peter Sellers",
		//	"Pres. Merkin Muffley":       "Peter Sellers",
		//	"Gen. Buck Turgidson":        "George C. Scott",
		//	"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
		//	`Maj. T.J. "King" Kong`:      "Slim Pickens",
		//},
		//Oscars: []string{
		//	"Best Actor (Nomin.)",
		//	"Best Adapted Screenplay (Nomin.)", "Best Director (Nomin.)",
		//	"Best Picture (Nomin.)",
		//},
	}

	marshal, err := Marshal(strangelove)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println(string(marshal))

	movie := Movie{}
	err = Unmarshal(marshal, &movie)
	if err != nil {
		fmt.Println("Unmarshal error", err)
		return
	}
	fmt.Println(movie)
}

type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }
func (lex *lexer) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}
	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct: // ((name value) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}
	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Bool:
		fmt.Fprintf(buf, "%t", v.Bool())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%e", v.Float())
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(buf, "%e", v.Complex())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	default: // float, complex, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

// Marshal encodes a Go value in S-expression form.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // get the first token
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}
