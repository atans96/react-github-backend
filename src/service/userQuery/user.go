package stargazersQuery

var QueryUser = `
    query ($login: String!) {
	  user(login: $login) {
		id
		login
		bio
		isGitHubStar
		location
		repositoryDiscussions(first: 1) {
		  pageInfo {
			endCursor
			hasNextPage
			startCursor
		  }
		  nodes {
			upvoteCount
			viewerHasUpvoted
			reactions(first: 5) {
			  pageInfo {
				endCursor
				hasNextPage
				startCursor
			  }
			  totalCount
			  nodes {
				content
				user {
				  id
				  login
				  bio
				  isGitHubStar
				  location
				}
			  }
			}
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
			category {
			  name
			}
		  }
		}
		repositoryDiscussionComments(first: 1) {
		  pageInfo {
			endCursor
			hasNextPage
			startCursor
		  }
		  nodes {
			upvoteCount
			viewerHasUpvoted
			discussion {
			  category {
				name
			  }
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
			reactions(first: 5) {
			  pageInfo {
				endCursor
				hasNextPage
				startCursor
			  }
			  totalCount
			  nodes {
				content
				user {
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
		starredRepositories(first: 2) {
		  pageInfo {
			endCursor
			hasNextPage
			startCursor
		  }
		  nodes {
			languages(first: 5) {
			  edges {
				node {
				  name
				}
				size
			  }
			}
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
			  }
			}
		  }
		}
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
			  participants(first: 5) {
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
			  participants(first: 5) {
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
