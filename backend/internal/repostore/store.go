package repostore

type Store interface {
	Close()
	User() UserRepository
	Actor() ActorRepository
	Film() FilmRepository
	FilmActor() FilmActorRepository
}
