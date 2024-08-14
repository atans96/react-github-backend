package stargazersQuery

var StarredRepositoriesMore = `
    query ($login: String!, $after: String!) {
	  user(login: $login) {
		id
		login
		bio
		isGitHubStar
		location
		starredRepositories(first: 50, after: $after) {
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
	  }
	}
`
