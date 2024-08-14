package stargazersQuery

var StarredRepositoriesStargazers = `
    query ($login: String!, $after: String!) {
	  user(login: $login) {
		id
		login
		bio
		isGitHubStar
		location
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
			stargazers(first: 50, after: $after) {
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
