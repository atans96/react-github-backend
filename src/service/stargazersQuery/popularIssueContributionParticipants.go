package stargazersQuery

var PopularIssueContributionParticipantsMore = `
    query ($owner: String!, $name: String!, $after: String!) {
	  repository(owner: $owner, name: $name) {
		stargazers(first: 1) {
		  pageInfo {
			endCursor
			hasNextPage
			startCursor
		  }
		  nodes {
			id
			login
			bio
			isGitHubStar
			location
			contributionsCollection {
			  popularIssueContribution {
				issue {
				  repository {
					defaultBranchRef {
                  		name
                    }
					name
					owner {
					  id
					  login
					}
					languages(first: 5) {
					  edges {
						node {
						  name
						}
						size
					  }
					}
				  }
				  participants(first: 50, after: $after) {
					pageInfo {
					  endCursor
					  hasNextPage
					  startCursor
					}
					nodes {
					  id
					  login
					  bio
					  isGitHubStar
					  location
					}
				  }
				}
			  }
			}
		  }
		}
	  }
	}
`
