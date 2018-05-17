package main

import "fmt"

type Store struct {
}

func (s *Store) Send(data string) {
	fmt.Println(data)
}
