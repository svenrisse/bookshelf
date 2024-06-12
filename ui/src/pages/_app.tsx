import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type AppType } from "next/app";
import { Inter } from "next/font/google";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import "~/styles/globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
});

const queryClient = new QueryClient();

const MyApp: AppType = ({ Component, pageProps }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <main className={`font-sans ${inter.variable}`}>
        <Component {...pageProps} />
      </main>
      <ReactQueryDevtools />
    </QueryClientProvider>
  );
};

export default MyApp;
