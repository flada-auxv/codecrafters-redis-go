package command

import "testing"

func TestOptParser_Parse(t *testing.T) {
	t.Run("When an Int option is set and passed", func(t *testing.T) {
		op := NewOptParser()
		ex := op.SetInt("EX")
		err := op.Parse([]string{"EX", "3"})
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if *ex != 3 {
			t.Errorf("got = %v, expected = %v", *ex, 3)
		}
	})

	t.Run("it return an error when it runs Parse twice", func(t *testing.T) {
		op := NewOptParser()
		ex := op.SetInt("EX")
		err := op.Parse([]string{"EX", "3"})
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if *ex != 3 {
			t.Errorf("got = %v, expected = %v", *ex, 3)
		}
		err2 := op.Parse([]string{"PX", "3000"})
		if err2.Error() != "already parsed" {
			t.Errorf("Expected error was not occurred")
		}
	})
}
