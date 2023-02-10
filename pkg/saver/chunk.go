package saver

type Chunk struct {
	Size   int
	Modulo int
}

func NewChunk(size int) Chunk {
	return Chunk{
		Size:   size / ChunkCount(),
		Modulo: size % ChunkCount(),
	}
}

func ChunkCount() int {
	return 5
}
