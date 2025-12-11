package main

func (c *commands) run(s *state, cmd command) error {
	return c.all[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.all[name] = f
	return nil
}
