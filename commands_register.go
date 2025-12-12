package main

func registerCommands(c *commands) {
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handleReset)
	c.register("users", handleUsers)
	c.register("agg", handleAgg)
	c.register("addfeed", handleAddFeed)
}
