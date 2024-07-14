package store

type Store interface {
	Close()
	User() UserRepository
	Actor() ActorRepository
	Film() FilmRepository
	FilmActor() FilmActorRepository
}
