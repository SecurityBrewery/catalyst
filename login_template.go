package catalyst

const temp = `
<html>
	<head>
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">\n
		<script src="//cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>\n
		<script src="//unpkg.com/alpinejs" defer></script>
		<script src="//unpkg.com/axios/dist/axios.min.js"></script>
		<script>
		</script>
	</head>
	<body style="background-color: #212121">
		<div class="container">
			<div class="card" style="margin-top: 120px">
				<div class="card-content">
					<span class="card-title">Login</span>
					
					%s

					<form action="/login">
						<div class="input-field col s12">
							<input id="username" name="username" type="text" class="validate">
							<label for="username">Username</label>
						</div>
						
						<div class="input-field col s12">
							<input id="password" name="password" type="password" class="validate">
							<label for="password">Password</label>
						</div>
						
						<button id="submit" type="submit" name="action" value="submit" class="waves-effect waves-light btn">Submit</button>
					</form>
				</div>
			</div>
		</div>
	</body>
</html>
`
