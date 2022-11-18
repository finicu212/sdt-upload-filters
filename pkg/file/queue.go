package file

type IFileQueue interface {
	AddToQueue(files []FileDetails)
	GetQueue() []FileDetails
	PopQueue() FileDetails
}

type FileQueue struct {
	Q []FileDetails
}

var _ IFileQueue = new(FileQueue) // Make sure FQ inherits IFQ

func (fq FileQueue) AddToQueue(files []FileDetails) {
	fq.Q = append(fq.Q, files...)
}

func (fq FileQueue) GetQueue() []FileDetails {
	return fq.Q
}

func (fq FileQueue) PopQueue() (next FileDetails) {
	next, fq.Q = fq.Q[0], fq.Q[1:]
	return next
}
