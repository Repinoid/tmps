package main

func main() {

	y := 0x77
	defer func() {
		y = 0x88
		_ = y
	}()

}
