package stargazersQuery

var RepositoryDiscussionsReactionsMore = `
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
			repositoryDiscussions(first: 1) {
			  pageInfo {
				endCursor
				hasNextPage
				startCursor
			  }
			  nodes {
				upvoteCount
				viewerHasUpvoted
				reactions(first: 50, after: $after) {
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
		  }
		}
	  }
	}
`
