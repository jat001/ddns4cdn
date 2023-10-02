package services

type Services interface {
	Run()
}

func Run(ctx *Services, id string, cha chan string) {
	defer func() {
		cha <- id
	}()
	(*ctx).Run()
}
