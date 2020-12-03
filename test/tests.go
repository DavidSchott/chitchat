/*Package test contains very basic sanity checks to aid in development and is not a substitute for actual tests */
package test

import (
	"fmt"

	"github.com/DavidSchott/chitchat/data"
)

func TestCreate() {
	cr1 := &data.ChatRoom{
		Title:       "Test Chat",
		Description: "Random test chat",
		Type:        "public",
		Password:    "",
		ID:          0,
	}
	cr1Dupe := &data.ChatRoom{
		Title:       "Test Chat",
		Description: "New description",
		Type:        "private",
		Password:    "asdasd",
		ID:          0,
	}
	cr2 := &data.ChatRoom{
		Title:       "Test Chat 2",
		Description: "Second test chat",
		Type:        "private",
		Password:    "3849583945",
		ID:          0,
	}
	cr3 := &data.ChatRoom{
		Title:       "chat 3",
		Description: "3 test chat",
		Type:        "hidden",
		Password:    "123",
		ID:          0,
	}
	data.CS.Add(cr1)
	data.CS.Add(cr2)
	data.CS.Add(cr3)
	if err := data.CS.Add(cr1Dupe); err != nil {
		fmt.Println((err.Error()))
	}
	fmt.Println("Created:", "\n", data.CS)
}

func TestRetrieve() {
	c1, _ := data.CS.Retrieve("public chat")
	if _, err := data.CS.Retrieve("not exist Chat"); err != nil {
		fmt.Println(err.Error())
	}
	c2, _ := data.CS.Retrieve("test chat 2")
	c3, _ := data.CS.Retrieve("2")
	c4, _ := data.CS.Retrieve("4")
	if _, err := data.CS.Retrieve("120"); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Retrieved:", "\n", c1, "\n", c2, "\n", c3, "\n", c4)
}
