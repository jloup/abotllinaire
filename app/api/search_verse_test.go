package api

import "testing"

func TestSearchVerse(t *testing.T) {

	t.Log(SearchVerse("soir", 3, "writer_out.txt", "writer_out_lower.txt"))

}
