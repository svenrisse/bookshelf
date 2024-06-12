# Instructions to work on UI

feel free to edit this file to take notes

## some preliminary considerations:

### do you want to use a component library?

- [shadcn/ui](https://ui.shadcn.com/) is considered one of the best on the market right now, if not **the** best
- [daisyui](https://daisyui.com/) is a well established alternative

### let's assume the stakeholders dont have any requirements for testing the ui.

- considering that, decide what should be tested regardless and what tools might be useful for that
- (vitest, playwright, cypress, jest, react-testing-library,...)

### deployment

- vercel is the go-to hosting provider for nextjs, react, etc.
- let's take 30 minutes together to deploy the ui to vercel and talk about the provided ci/cd

## data fetching

- [React-Query](https://tanstack.com/query/latest) is the go-to solution for data-fetching in react.
  Spend 5-10 Minutes reading the docs to understand why it should be used in this project.
- I took the liberty to init React-Query, the eslint-plugin and the React-Query-Devtools for you in [863b9a361ff1](https://github.com/svenrisse/bookshelf/commit/863b9a361ff1cf3465ec72e37d21cc53ac810369)

## first issue

- [#11](https://github.com/svenrisse/bookshelf/issues/11) should be a good place to start.
- I'll outline the steps for this issue, to get you feeling comfortable working on this project

- [ ] assign yourself this issue by moving it from "UI-Backlog" to "In progress" on the [github project](https://github.com/users/svenrisse/projects/3/views/1)
- [ ] create a branch with a fitting name, using the github index as a prefix, like 11-bookDetailsPage
- [ ] push the local branch to github, to get preview deployments and working ci/cd
- [ ] we want the url for this page to be /book/show/{id}, use nextjs [slugs](https://nextjs.org/docs/pages/building-your-application/routing/dynamic-routes) for that
- [ ] one way of creating user interfaces is using tones of grey at the start and implement colors later, feel free to give this a try
- [ ] don't worry about data fetching for now and build responsive components with dummy data (remember to use the provided goodreads link for reference)
- [ ] for now just render 5 star icons, we'll worry about correctly coloring them later
- [ ] looking at the goodreads example, you'll see a lot of links. think about more /pages/... to implement and when you would need slugs
- [ ] it's time to fetch the data from the API, lets look at this step in a dedicated section

## data fetching with React-Query

![image](https://github.com/svenrisse/bookshelf/assets/89209935/48cb2147-2259-4b34-8f62-abf54fb85d49)

