package stargazersQuery

var PopularIssueContributionParticipantsMore = `
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
		  popularIssueContribution {
			issue {
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
`
