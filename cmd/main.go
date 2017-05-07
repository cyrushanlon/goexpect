package main

import (
	"os"
	"time"

	"github.com/cyrushanlon/expect"
	"github.com/yuin/gopher-lua"
)

var (
	p expect.Process
)

func spawnP(L *lua.LState) int { //*
	cmd := L.ToString(1) // get first (1) function argument and convert to int
	args := L.ToString(2)

	p = expect.Process{}
	p.Timeout = 5
	p.Start(cmd, args)

	//ln := lua.LBool(true) // make calculation and cast to LNumber
	//L.Push(ln)            // Push it to the stack
	return 0 // Notify that we pushed one value to the stack
}

func expectP(L *lua.LState) int {
	check := L.ToString(1)

	ln := lua.LBool(p.Expect(check, true)) // make calculation and cast to LNumber
	L.Push(ln)                             // Push it to the stack

	return 1 // Notify that we pushed one value to the stack
}

func sendP(L *lua.LState) int {
	sent := L.ToString(1)

	p.SendInput(sent)

	return 0 // Notify that we pushed one value to the stack
}

func sleep(L *lua.LState) int {
	ms := L.ToInt(1)

	time.Sleep(time.Duration(ms) * time.Millisecond)

	return 0 // Notify that we pushed one value to the stack
}

func exit(L *lua.LState) int {
	p.Close()

	return 0 // Notify that we pushed one value to the stack
}

func main() {

	argv := os.Args[1:] //remove first arg which is the path to the exe

	L := lua.NewState()
	defer L.Close()

	argvLuaTable := L.NewTable()
	for k, v := range argv {
		argvLuaTable.Insert(k, lua.LString(v))
	}

	L.SetGlobal("argv", argvLuaTable)

	//register the functions from Go to Lua
	L.SetGlobal("spawn", L.NewFunction(spawnP))   // Register our function in Lua
	L.SetGlobal("send", L.NewFunction(sendP))     // Register our function in Lua
	L.SetGlobal("expect", L.NewFunction(expectP)) // Register our function in Lua
	L.SetGlobal("sleep", L.NewFunction(sleep))    // Register our function in Lua
	L.SetGlobal("exit", L.NewFunction(exit))      // Register our function in Lua

	if err := L.DoString(`
	
	for k, v in pairs(argv) do
		print(v) 
	end

	spawn("cmd", "/C ping 8.8.8.8")

	sleep(500)

	send("hello")

	if expect("approximate round trip times in milli-seconds:") then
		print("woop de doo")
	else
		print("oh no de oh")
	end

	exit()

	`); err != nil {
		panic(err)
	}
}
