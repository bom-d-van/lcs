package lcs

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"sync"
)

type Prose interface {
	New() Prose
	Len() int
	Word(i int) Word
	Slice(i, j int) Prose
	Append(op Prose)
	AppendWord(w Word)
	WrapDel()
	WrapIns()
}

type Word interface {
	IsEqual(ow Word) bool
	IndexIn(p Prose, i int) int
}

type article struct {
	id    int
	terms []*term
}

var puncs = regexp.MustCompile("[,!;:\"?\\.]$")

func newArticle(ori io.Reader) *article {
	a := article{id: genSpId()}
	scanner := bufio.NewScanner(ori)
	scanner.Split(scanWords)

	for scanner.Scan() {
		w := scanner.Text()
		a.terms = append(a.terms, distillTerms(w)...)
	}
	return &a
}

func distillTerms(w string) (r []*term) {
	if i := strings.Index(w, "\n"); i != -1 {
		if i == 0 {
			r = append(r, newTerm(string(w[i])))
		} else {
			r = append(r, distillTerms(w[:i])...)
			r = append(r, newTerm(string(w[i])))
		}

		if i+1 < len(w) {
			r = append(r, distillTerms(w[i+1:])...)
		}

		return
	}

	terms := []*term{}

	for puncs.MatchString(w) {
		terms = append(terms, newTerm(string(w[len(w)-1])))
		w = w[:len(w)-1]
	}

	if w != "" {
		terms = append(terms, newTerm(w))
	}

	for i, _ := range terms {
		r = append(r, terms[len(terms)-i-1])
	}

	return
}

func (a *article) String() (str string) {
	for i, term := range a.terms {
		if !term.isPuncs() {
			if i == 0 {
				if term.isLF() || a.terms[i+1].isLF() {
					goto appendTerm
				}
			} else {
				if a.terms[i-1].isLF() || a.terms[i+1].isLF() {
					goto appendTerm
				}

				if !a.terms[i-1].isOpenMarks() && !term.isCloseMarks() {
					str += " "
				}
			}
		}
	appendTerm:
		str += term.string
	}

	return
}

func (t *term) isPuncs() bool {
	return puncs.MatchString(t.string)
}

func (t term) isOpenMarks() bool {
	return t.string == "[" || t.string == "("
}

func (t term) isLF() bool {
	return t.string == "\n"
}

func (t term) isCloseMarks() bool {
	return t.string == "]" || t.string == ")"
}

func (a article) New() Prose {
	return &article{id: genSpId()}
}

func (a article) Len() int {
	return len(a.terms)
}

func (a article) Word(i int) Word {
	return a.terms[i]
}

func (a article) Slice(i, j int) Prose {
	na := new(article)
	na.id = a.id
	na.terms = a.terms[i:j]
	return na
}

func (a *article) Append(op Prose) {
	oa := op.(*article)
	// a.origin += oa.origin
	a.terms = append(a.terms, oa.terms...)
}

func (a *article) Prepend(op Prose) {
	oa := op.(*article)
	// a.origin = oa.origin + a.origin
	a.terms = append(oa.terms, a.terms...)
}

func (a *article) AppendWord(w Word) {
	t := w.(*term)
	// a.origin += t.string
	a.terms = append(a.terms, t)
}

func (a *article) PrependWord(w Word) {
	t := w.(*term)
	// a.origin = t.string + a.origin
	a.terms = append([]*term{t}, a.terms...)
}

func (a *article) WrapDel() {
	a.AppendWord(newTerm("]"))
	// a.origin = "[" + a.origin
	a.terms = append([]*term{newTerm("[")}, a.terms...)
}

func (a *article) WrapIns() {
	a.AppendWord(newTerm(")"))
	// a.origin = "(" + a.origin
	a.terms = append([]*term{newTerm("(")}, a.terms...)
}

type term struct {
	string
	pos map[int]int
}

func newTerm(cont string) *term {
	return &term{string: cont, pos: map[int]int{}}
}

func (t term) IsEqual(ow Word) bool {
	return t.string == ow.(*term).string
}

func (t *term) IndexIn(p Prose, i int) int {
	if i == -1 {
		return t.pos[p.(*article).id]
	}
	t.pos[p.(*article).id] = i
	return i
}

func (t *term) String() string {
	return t.string
}

type stringProse struct {
	// string
	words []*byteWord
	id    int
}

var (
	spId     = 0
	spIdLock = sync.Mutex{}
)

func genSpId() int {
	spIdLock.Lock()
	defer spIdLock.Unlock()
	spId += 1

	return spId
}

func newStringProse(content string) (sp *stringProse) {
	sp = new(stringProse)
	sp.id = genSpId()
	for _, b := range content {
		w := newByteWord(byte(b))
		// w.indexIn(sp, i)
		sp.words = append(sp.words, w)
	}

	return
}

func (s stringProse) Len() int {
	return len(s.words)
}

func (s stringProse) Word(i int) Word {
	return s.words[i]
}

func (s stringProse) Slice(i, j int) Prose {
	p := newStringProse(s.SliceString(i, j))
	p.id = s.id
	return p
}

func (s stringProse) SliceString(i, j int) (str string) {
	for _, w := range s.words[i:j] {
		str += string(w.byte)
	}
	return
}

func (s *stringProse) Append(op Prose) {
	s.words = append(s.words, op.(*stringProse).words...)
}

func (s *stringProse) Prepend(op Prose) {
	s.words = append(op.(*stringProse).words, s.words...)
}

func (s *stringProse) PrependWord(w Word) {
	s.words = append([]*byteWord{w.(*byteWord)}, s.words...)
}

func (s *stringProse) AppendWord(w Word) {
	s.words = append(s.words, w.(*byteWord))
}

func (s stringProse) New() Prose {
	return &stringProse{[]*byteWord{}, genSpId()}
}

func (s *stringProse) WrapDel() {
	s.PrependWord(newByteWord('['))
	s.AppendWord(newByteWord(']'))
}

func (s *stringProse) WrapIns() {
	s.PrependWord(newByteWord('('))
	s.AppendWord(newByteWord(')'))
}

func (s stringProse) String() (str string) {
	for _, w := range s.words {
		str += string(w.byte)
	}
	return
}

type byteWord struct {
	byte byte
	pos  map[int]int
}

func newByteWord(b byte) (bw *byteWord) {
	bw = new(byteWord)
	bw.byte = b
	bw.pos = map[int]int{}
	return
}

func (b byteWord) IsEqual(ow Word) bool {
	return byte(b.byte) == byte(ow.(*byteWord).byte)
}

func (b *byteWord) IndexIn(p Prose, i int) int {
	if i == -1 {
		return b.pos[p.(*stringProse).id]
	}
	b.pos[p.(*stringProse).id] = i
	return i
}

func (b byteWord) String() string {
	return string(b.byte)
}

func LCS(a, b Prose) Prose {
	aLen := a.Len()
	bLen := b.Len()
	lengths := make([][]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		lengths[i] = make([]int, bLen+1)
	}

	// row 0 and column 0 are initialized to 0 already
	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			if a.Word(i).IsEqual(b.Word(j)) {
				lengths[i+1][j+1] = lengths[i][j] + 1
			} else if lengths[i+1][j] > lengths[i][j+1] {
				lengths[i+1][j+1] = lengths[i+1][j]
			} else {
				lengths[i+1][j+1] = lengths[i][j+1]
			}
		}
	}

	// for _, r := range lengths {
	// 	fmt.Println(r)
	// }

	// read the substring out from the matrix
	s := a.New()
	for x, y := aLen, bLen; x != 0 && y != 0; {
		if lengths[x][y] == lengths[x-1][y] {
			x--
		} else if lengths[x][y] == lengths[x][y-1] {
			y--
		} else {
			w := a.Word(x - 1)
			w.IndexIn(a, x-1)
			w.IndexIn(b, y-1)
			s.AppendWord(w)
			x--
			y--
		}
	}

	// reverse string
	r := a.New()
	for i := 0; i < s.Len(); i++ {
		r.AppendWord(s.Word(s.Len() - 1 - i))
	}
	return r
}

func Diff(ori, edit, lcs Prose) Prose {
	// println("---->")
	// fmt.Println(ori, edit, lcs)
	if lcs.Len() == 0 {
		ori.WrapDel()
		edit.WrapIns()
		ori.Append(edit)
		return ori
	}

	diff := ori.New()
	lastIOri := -1
	lastIEdit := -1
	lenOri := ori.Len()
	lenEdit := edit.Len()
	lenLCS := lcs.Len()
	for i := 0; i < lenLCS; i++ {
		w := lcs.Word(i)
		oi := w.IndexIn(ori, -1)
		if lenOri > i && oi != i {
			// deletion
			// fmt.Printf("--> %+v\n", w.(*term).pos)
			// println(w.(*term).string)
			// println(lastIOri+1, oi)
			// fmt.Println(ori.(*article).id)
			del := ori.Slice(lastIOri+1, oi)
			if del.Len() > 0 {
				del.WrapDel()
				diff.Append(del)
			}
		}
		lastIOri = oi
		ei := w.IndexIn(edit, -1)
		if lenEdit > i && ei != i {
			// insertion
			ins := edit.Slice(lastIEdit+1, ei)
			if ins.Len() > 0 {
				ins.WrapIns()
				diff.Append(ins)
			}
		}
		diff.AppendWord(w)
		lastIEdit = ei
	}

	// deletion and insertion after last word of ori and edit
	lastw := lcs.Word(lcs.Len() - 1)
	lastIOri = lastw.IndexIn(ori, -1)
	if lenOri > lastIOri+1 {
		del := ori.Slice(lastIOri+1, lenOri)
		if del.Len() > 0 {
			del.WrapDel()
			diff.Append(del)
		}
	}
	lastIEdit = lastw.IndexIn(edit, -1)
	if lenEdit > lastIEdit+1 {
		ins := edit.Slice(lastIEdit+1, lenEdit)
		if ins.Len() > 0 {
			ins.WrapIns()
			diff.Append(ins)
		}
	}

	return diff
}
