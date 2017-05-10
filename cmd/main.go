package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	expect "github.com/cyrushanlon/goexpect"
	"github.com/yuin/gopher-lua"
)

var (
	p        expect.Process
	exitCode int
)

func spawnP(L *lua.LState) int { //*

	cmd := L.ToString(1) // get first (1) function argument and convert to int
	args := L.ToString(2)

	p = expect.Process{}
	p.Timeout = 5

	err := p.Start(cmd, args)
	if err != nil {
		fmt.Println(err)
	}

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
	time.Sleep(50 * time.Millisecond)

	return 0 // Notify that we pushed one value to the stack
}

func sleep(L *lua.LState) int {
	ms := L.ToInt(1)

	time.Sleep(time.Duration(ms) * time.Millisecond)

	return 0 // Notify that we pushed one value to the stack
}

func exit(L *lua.LState) int {
	exitCode = L.ToInt(1) // get first (1) function argument and convert to int
	time.Sleep(500 * time.Millisecond)

	p.Close()

	return 0 // Notify that we pushed one value to the stack
}

func setTimeout(L *lua.LState) int { //*

	i := L.ToInt(1) // get first (1) function argument and convert to int

	p.Timeout = i

	return 0 // Notify that we pushed one value to the stack
}

func main() {

	p = expect.Process{}

	//for very low powered machines nothing would happen in the goroutines as soon as they are put to sleep
	if runtime.GOMAXPROCS(0) < 3 {
		runtime.GOMAXPROCS(3)
	}

	argv := os.Args

	if len(argv) < 2 {
		fmt.Println("no path to script")
		fmt.Println("exitCode error!")
		os.Exit(1)
		//panic("no path to script")
	}

	L := lua.NewState()
	defer L.Close()

	argvLuaTable := L.NewTable()
	//remove first arg which is the path to the exe
	//remove second arg as it is the path to the script
	for k, v := range argv[2:] {
		argvLuaTable.Insert(k, lua.LString(v))
	}

	L.SetGlobal("argv", argvLuaTable)

	//register the functions from Go to Lua
	L.SetGlobal("spawn", L.NewFunction(spawnP))       // Register our function in Lua
	L.SetGlobal("send", L.NewFunction(sendP))         // Register our function in Lua
	L.SetGlobal("expect", L.NewFunction(expectP))     // Register our function in Lua
	L.SetGlobal("sleep", L.NewFunction(sleep))        // Register our function in Lua
	L.SetGlobal("exit", L.NewFunction(exit))          // Register our function in Lua
	L.SetGlobal("timeout", L.NewFunction(setTimeout)) // Register our function in Lua

	if err := L.DoFile(argv[1]); err != nil {
		fmt.Println(err)
		fmt.Println("exitCode error!")
		os.Exit(1)
		//panic(err)
	}

	if exitCode != 0 {
		fmt.Println("exitCode error!")
	}
	os.Exit(exitCode)
}
