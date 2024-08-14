package stargazersQuery

var RepositoryDiscussionCommentsReactionsMore = `
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
			  }
			}
		  }
		}
	  }
	}
`
