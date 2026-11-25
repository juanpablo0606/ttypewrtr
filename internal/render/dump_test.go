package render

import "testing"

func TestDump(t *testing.T) {
	src := []byte("# Título **negrita** y *cursiva*\n\n" +
		"> una cita\n\n" +
		"- un ítem\n\n" +
		"**a*b** \\*escapado\\* `**código**`\n")
	t.Log("\n" + Dump(src))
}
