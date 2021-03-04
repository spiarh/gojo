package core

// func SetFacts(flagSet *pflag.FlagSet, facts []*Fact, sources []Source) error {
// 	for _, fact := range facts {
// 		if fact.Source == "" {
// 			continue
// 		}
// 		for _, src := range sources {
// 			if fact.Source == src.Name {
// 				repo, err := provider.New(flagSet, src)
// 				if err != nil {
// 					return err
// 				}

// 				if fact.Value, err = repo.GetLatest(); err != nil {
// 					return err
// 				}

// 				if fact.Kind == VersionFactKind {
// 					fact.Value = util.SanitizeVersion(fact.Value)
// 				}
// 			}
// 		}
// 		if fact.Value == "" {
// 			log.Fatal().Msgf("no value found for fact with name: %s", fact.Name)
// 		}
// 	}
// 	return nil
// }
