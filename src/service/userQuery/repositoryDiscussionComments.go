package stargazersQuery

var RepositoryDiscussionComments = `
    query ($login: String!, $after: String!) {
	  user(login: $login) {
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
			repositoryDiscussionComments(first: 50, after: $after) {
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
		  }
		}
	  }
	}
`
