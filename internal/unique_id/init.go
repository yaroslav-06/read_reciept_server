package uniqueid

import (
	"fmt"
	"math/rand"
)

type Generator struct {
	used_ids map[string]bool
}

func NewGenerator() *Generator {
	gr := &Generator{}
	gr.used_ids = make(map[string]bool)
	return gr
}

func (gr *Generator) DoesIdExists(id string) bool {
	vl, exist := gr.used_ids[id]
	return exist && vl
}

func (gr *Generator) GetNewId() string {
	for {
		new_id := fmt.Sprint(rand.Int63())
		if !gr.DoesIdExists(new_id) {
			gr.used_ids[new_id] = true
			return new_id
		}
	}
}
