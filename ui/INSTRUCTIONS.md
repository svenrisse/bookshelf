# Instructions to work on UI

feel free to edit this file to take notes

## some preliminary considerations:

#### do you want to use a component library?

- [shadcn/ui](https://ui.shadcn.com/) is considered one of the best on the market right now, if not **the** best
- [daisyui](https://daisyui.com/) is a well established alternative

#### let's assume the stakeholders dont have any requirements for testing the ui.

- considering that, decide what should be tested regardless and what tools might be useful for that
- (vitest, playwright, cypress, jest, react-testing-library,...)

#### deployment

- vercel is the go-to hosting provider for nextjs, react, etc.
- let's take 30 minutes together to deploy the ui to vercel and talk about the provided ci/cd

## data fetching

- [React-Query](https://tanstack.com/query/latest) is the go-to solution for data-fetching in react.
  Spend 5-10 Minutes reading the docs to understand why it should be used in this project.
- I took the liberty to init React-Query, the eslint-plugin and the React-Query-Devtools for you in [863b9a361ff1](https://github.com/svenrisse/bookshelf/commit/863b9a361ff1cf3465ec72e37d21cc53ac810369)
