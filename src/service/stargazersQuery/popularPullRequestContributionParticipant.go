package stargazersQuery

var PopularPullRequestContributionParticipantsMore = `
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
			  totalIssueContributions
			  totalPullRequestContributions
			  totalPullRequestReviewContributions
			  pullRequestReviewContributionsByRepository {
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
			  }
			  pullRequestContributionsByRepository {
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
			  }
			  popularPullRequestContribution {
				pullRequest {
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
	  }
	}
`
