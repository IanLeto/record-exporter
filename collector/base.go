package collector

type DataSource interface {
	GetData() error
}
