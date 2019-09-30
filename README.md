# community

## External contributions

You might have access here to have an idea of how the project is progressing,
or to perhaps help provide insight on some questions I might have throughout the year.

Please keep in mind that this project is being used to get "academic credit",
and I need to be very careful to not misrepresent the work I have done!

I ask that you do not push anything to this repository, and to refrain from making code contributions.

**Question to Don (Honours project coordinator)**
> Is source for honours projects required to be closed source until the end of the project?
> i.e can it be an open source project?
>
> I'm aware we own the code and can do whatever after project submission, but I am not sure what the general policy is before submission of the project.

**Answer**
> You can publish the code as open source when the project is running, but since the project must be your own work you can't incorporate contributions from other people. Or if you do so, you need to make it crystal clear what is your work and what is the work of other people. If it's an active open source project with a community of contributors, that's probably not possible.

## Docs

Build using `go build`.

Copy `config.yaml.example` to `config.yaml`.

Start with `config=config.yaml ./community`

## OAuth

- Send the user to `https://forum.mtasa.com/oauth/authorize/?client_id={CLIENT_ID}&response_type=code&redirect_uri=http://localhost:8080`
