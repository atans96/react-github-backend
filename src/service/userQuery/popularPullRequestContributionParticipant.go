package stargazersQuery

var PopularPullRequestContributionParticipantsMore = `
    query ($login: String!, $after: String!) {
	  user(login: $login) {
		id
		login
		bio
		isGitHubStar
		location
		contributionsCollection {
		  totalIssueContributions
		  totalPullRequestContributions
		  totalPullRequestReviewContributions
		  pullRequestReviewContributionsByRepository {
			repository {
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
		  }
		  pullRequestContributionsByRepository {
			repository {
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
		  }
		  popularPullRequestContribution {
			pullRequest {
			  repository {
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
				nodes {
				  id
				  login
				  bio
				  isGitHubStar
				  location
				}
				pageInfo {
				  endCursor
				  hasNextPage
				  startCursor
				}
			  }
			}
		  }
		}
	  }
	}
`
