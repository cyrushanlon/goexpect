package main

import (
	"os"
	"time"

	"log"

	"github.com/cyrushanlon/goexpect"
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

	err := p.Start(cmd, args)
	if err != nil {
		log.Println(err)
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
	time.Sleep(15 * time.Millisecond)

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

func setTimeout(L *lua.LState) int { //*

	i := L.ToInt(1) // get first (1) function argument and convert to int

	p.Timeout = i

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
	L.SetGlobal("spawn", L.NewFunction(spawnP))       // Register our function in Lua
	L.SetGlobal("send", L.NewFunction(sendP))         // Register our function in Lua
	L.SetGlobal("expect", L.NewFunction(expectP))     // Register our function in Lua
	L.SetGlobal("sleep", L.NewFunction(sleep))        // Register our function in Lua
	L.SetGlobal("exit", L.NewFunction(exit))          // Register our function in Lua
	L.SetGlobal("timeout", L.NewFunction(setTimeout)) // Register our function in Lua

	if err := L.DoString(`
	--first we sort out the vars
	timeout(5)
	ip = argv[0]
	username =  argv[1]
	password = argv[2]
	macs = {}
	for i = 3, #argv do
    	macs[i-2] = argv[i]
  	end

	--then we create the telnet instance
	spawn("telnet", ip)
	
	--login
	if not expect(":") then
		exit()
		return
	end
	send(username)

	if not expect(":") then
		exit()
		return
	end
	send(password)

	if not expect(">") then
		exit()
		return
	end
	--we are now logged in

	--reboot in case it all goes completly wrong at some point in the future
	send("reboot 8")
	if not expect("reboot time 8 minutes") and not expect("ok") then
		exit()
		return
	end

	--mac address filter
	for k, v in pairs(macs) do
		send ("macfilt ".. k .. " mac ".. v)
		if not expect("ok") then
			exit()
			return
		end
	end
	
	--enable filter
	send("eth 0 macfilt on")
	if not expect("ok") then
		exit()
		return
	end
	send("exit")
	
	exit()

	`); err != nil {
		panic(err)
	}
}
