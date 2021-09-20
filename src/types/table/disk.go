package table

import (
	"fmt"
	"os"
	"strings"
)

func (t *Table) ToDisk(path string) error {
	/*
		one file per column
		create one file every a 10k row
	*/
	a := []string{}
	ty := []byte{}
	e := make(chan error)
	s := make(chan bool)
	for k, v := range t.Columns {
		go func() {
			err := v.ToDisk(path)
			if err != nil {
				e <- err
			} else {
				s <- true
			}

		}()
		a = append(a, k)
		ty = append(t, v.Dtype)
	}
	go func() {

		f, err := os.Create(path + string(os.PathSeparator) + "meta.json")
		if err != nil {
			e <- err
		}

		b := "['" + strings.Join(a, "','") + "']"
		_, err = f.Write([]byte(fmt.Sprintf("{'permission':%v,'Column_names':%v,'Dtypes':%v}", t.permission, b, string(ty))))

		if err != nil {
			e <- err
		}

	}()
	var c int = 0
	d := make(chan bool)

	go func() {
		for {
			_ = <-s
			c += 1
			if c == len(t.Columns+1) {
				d <- true
			}
		}
	}()

	select {
	case err := <-e:
		return err
	case <-d:
		return nil
	}
	return nil
}
