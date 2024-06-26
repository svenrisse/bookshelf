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

this is a simple get request with react-query using fetch under the hood

![image](https://github.com/svenrisse/bookshelf/assets/89209935/48cb2147-2259-4b34-8f62-abf54fb85d49)

- using an inline queryFn is most often the right way of doing things
- consider moving the useQuery call into its own function/file to make it a custom hook, this makes testing and reusing it way easier
- take a look at all the vars and methods useQuery returns and what purpose they might serve

---

![image](https://res.cloudinary.com/practicaldev/image/fetch/s--h_E-Gwvm--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_800/https://cdn.hashnode.com/res/hashnode/image/upload/v1643349189447/XLfQFf761.png)

- take note of the react-query-devtools while implementing your own api call, it makes debuging and understanding what might be going wrong a breeze
- you can find it on your localhost instance in one corner

## loading state

- think about how you want to show a loading state to the user
- you might want to use a [skeleton](https://ui.shadcn.com/docs/components/skeleton)
- use one of the booleans returned from useQuery to get the loading state
- tip: for development it's useful to implement the loading skeletons for the !loading state (you don't have to keep refreshing)
