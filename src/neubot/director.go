package main

func DirectorRun(neubot_home string, nettest_name string) (error) {
	spec, err := SpecLoad(neubot_home, nettest_name)
	if err != nil {
		return err
	}
	return SpecRunSync(spec)
}
