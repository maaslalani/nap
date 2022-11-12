package main

import "fmt"

func foo() string {
	for i := 0; i < 10; i++ {
		for i := 0; i < 10; i++ {
			for i := 0; i < 10; i++ {
				for i := 0; i < 10; i++ {
					for i := 0; i < 10; i++ {
						for i := 0; i < 10; i++ {
							for i := 0; i < 10; i++ {
								for i := 0; i < 10; i++ {
									for i := 0; i < 10; i++ {
										for i := 0; i < 10; i++ {
											fmt.Println("Please, no.")
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "foo"
}
