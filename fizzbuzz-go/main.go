package main

func main() {
	for i := 0; i < 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			print("FizzBuzz")
		} else if i%3 == 0 {
			print("Fizz")
		} else if i%5 == 0 {
			print("Buzz")
		} else {
			print(i)
		}
		print(" ")
	}
}
