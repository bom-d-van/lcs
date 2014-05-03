package lcs

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"sync"
)

type prose interface {
	new() prose
	len() int
	word(i int) word
	slice(i, j int) prose
	append(op prose)
	// prepend(op prose)
	appendWord(w word)
	// prependWord(w word)
	wrapDel()
	wrapIns()
}

type word interface {
	isEqual(ow word) bool
	indexIn(p prose, i int) int
}

type article struct {
	id int
	// origin string
	terms []*term
}

var (
	// puncs = ",!;:\"?.\n"
	puncs = regexp.MustCompile("[,!;:\"?\\.]$")
)

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

func (a article) new() prose {
	return &article{id: genSpId()}
}

func (a article) len() int {
	return len(a.terms)
}

func (a article) word(i int) word {
	return a.terms[i]
}

func (a article) slice(i, j int) prose {
	na := new(article)
	na.id = a.id
	na.terms = a.terms[i:j]
	return na
}

func (a *article) append(op prose) {
	oa := op.(*article)
	// a.origin += oa.origin
	a.terms = append(a.terms, oa.terms...)
}

func (a *article) prepend(op prose) {
	oa := op.(*article)
	// a.origin = oa.origin + a.origin
	a.terms = append(oa.terms, a.terms...)
}

func (a *article) appendWord(w word) {
	t := w.(*term)
	// a.origin += t.string
	a.terms = append(a.terms, t)
}

func (a *article) prependWord(w word) {
	t := w.(*term)
	// a.origin = t.string + a.origin
	a.terms = append([]*term{t}, a.terms...)
}

func (a *article) wrapDel() {
	a.appendWord(newTerm("]"))
	// a.origin = "[" + a.origin
	a.terms = append([]*term{newTerm("[")}, a.terms...)
}

func (a *article) wrapIns() {
	a.appendWord(newTerm(")"))
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

func (t term) isEqual(ow word) bool {
	return t.string == ow.(*term).string
}

func (t *term) indexIn(p prose, i int) int {
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

func (s stringProse) len() int {
	// return len(s.string)
	return len(s.words)
}

func (s stringProse) word(i int) word {
	return s.words[i]
}

func (s stringProse) slice(i, j int) prose {
	// sp := stringProse(s[i:j])
	// return &stringProse{s.string[i:j], genSpId()}
	// p.words = make([]byteWord, j-i)
	// copy(p.words, s.words[i:j])
	p := newStringProse(s.sliceString(i, j))
	p.id = s.id
	return p
}

func (s stringProse) sliceString(i, j int) (str string) {
	for _, w := range s.words[i:j] {
		str += string(w.byte)
	}
	return
}

func (s *stringProse) append(op prose) {
	// s.string += op.(*stringProse).string
	s.words = append(s.words, op.(*stringProse).words...)
}

func (s *stringProse) prepend(op prose) {
	// *s = stringProse(string(*op.(*stringProse)) + string(*s))
	// s.string = op.(*stringProse).string + s.string
	s.words = append(op.(*stringProse).words, s.words...)
}

func (s *stringProse) prependWord(w word) {
	// *s = stringProse(string(byte(w.(byteWord).byte)) + string(*s))
	// s.string = string(w.(*byteWord).byte) + s.string
	s.words = append([]*byteWord{w.(*byteWord)}, s.words...)
}

func (s *stringProse) appendWord(w word) {
	s.words = append(s.words, w.(*byteWord))
}

func (s stringProse) new() prose {
	// np := stringProse("")
	return &stringProse{[]*byteWord{}, genSpId()}
}

func (s *stringProse) wrapDel() {
	// s.string = "[" + s.string + "]"
	s.prependWord(newByteWord('['))
	s.appendWord(newByteWord(']'))
}

func (s *stringProse) wrapIns() {
	// s.string = "(" + s.string + ")"
	s.prependWord(newByteWord('('))
	s.appendWord(newByteWord(')'))
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

func (b byteWord) isEqual(ow word) bool {
	return byte(b.byte) == byte(ow.(*byteWord).byte)
}

func (b *byteWord) indexIn(p prose, i int) int {
	if i == -1 {
		return b.pos[p.(*stringProse).id]
	}
	b.pos[p.(*stringProse).id] = i
	return i
}

func (b byteWord) String() string {
	return string(b.byte)
}

func LCS(a, b prose) prose {
	aLen := a.len()
	bLen := b.len()
	lengths := make([][]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		lengths[i] = make([]int, bLen+1)
	}

	// row 0 and column 0 are initialized to 0 already
	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			if a.word(i).isEqual(b.word(j)) {
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
	s := a.new()
	for x, y := aLen, bLen; x != 0 && y != 0; {
		if lengths[x][y] == lengths[x-1][y] {
			x--
		} else if lengths[x][y] == lengths[x][y-1] {
			y--
		} else {
			w := a.word(x - 1)
			w.indexIn(a, x-1)
			w.indexIn(b, y-1)
			s.appendWord(w)
			x--
			y--
		}
	}

	// reverse string
	r := a.new()
	for i := 0; i < s.len(); i++ {
		r.appendWord(s.word(s.len() - 1 - i))
	}
	return r
}

func Diff(ori, edit, lcs prose) prose {
	// println("---->")
	// fmt.Println(ori, edit, lcs)
	if lcs.len() == 0 {
		ori.wrapDel()
		edit.wrapIns()
		ori.append(edit)
		return ori
	}

	diff := ori.new()
	lastIOri := -1
	lastIEdit := -1
	lenOri := ori.len()
	lenEdit := edit.len()
	lenLCS := lcs.len()
	for i := 0; i < lenLCS; i++ {
		w := lcs.word(i)
		oi := w.indexIn(ori, -1)
		if lenOri > i && oi != i {
			// deletion
			// fmt.Printf("--> %+v\n", w.(*term).pos)
			// println(w.(*term).string)
			// println(lastIOri+1, oi)
			// fmt.Println(ori.(*article).id)
			del := ori.slice(lastIOri+1, oi)
			if del.len() > 0 {
				del.wrapDel()
				diff.append(del)
			}
		}
		lastIOri = oi
		ei := w.indexIn(edit, -1)
		if lenEdit > i && ei != i {
			// insertion
			ins := edit.slice(lastIEdit+1, ei)
			if ins.len() > 0 {
				ins.wrapIns()
				diff.append(ins)
			}
		}
		diff.appendWord(w)
		lastIEdit = ei
	}

	// deletion and insertion after last word of ori and edit
	lastw := lcs.word(lcs.len() - 1)
	lastIOri = lastw.indexIn(ori, -1)
	if lenOri > lastIOri+1 {
		del := ori.slice(lastIOri+1, lenOri)
		if del.len() > 0 {
			del.wrapDel()
			diff.append(del)
		}
	}
	lastIEdit = lastw.indexIn(edit, -1)
	if lenEdit > lastIEdit+1 {
		ins := edit.slice(lastIEdit+1, lenEdit)
		if ins.len() > 0 {
			ins.wrapIns()
			diff.append(ins)
		}
	}

	return diff
}
