import { ErrorTooltip } from "@/component/errorTooltip";
import AuthService from "@/service/auth.service";
import { Button } from "@/vendor/shadcn/components/ui/button";
import { Card, CardContent } from "@/vendor/shadcn/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem
} from "@/vendor/shadcn/components/ui/form";
import { Input } from "@/vendor/shadcn/components/ui/input";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

// FIXME: this is a bit ugly but it currently works, should create a nicer UI in the future
function wrongUserOrPasswordTooltip(show: boolean) {
    if (!show) return <></>
    return (
        <div>
            <ErrorTooltip 
                message="Wrong username or password, please try again."
            />
        </div>
    );
}


export function LoginPage() {
  const form = useForm();
  // TODO: create navigator that already computes prefix path
  const navigate = useNavigate();

  const [showWrongCredentialsTooltip, setShowWrongCredentialsTooltip] = useState<boolean>(false);

  if(AuthService.isAuthenticated()) {
    AuthService.logout();
  }

  return (
    <div className="flex h-screen items-center justify-center">
      <Card className="w-100 space-y-8">
        <CardContent className="py-5 px-10 bg-zinc-850">
          <h1 className="text-2xl font-bold pb-8">Enter your control panel</h1>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit((event) => {
                console.log("onsubmit", event);
                const { email, password } = event;

                AuthService.tryLogin({ email, password })
                  .then(() => {
                    console.log('resolved')
                    navigate(`${import.meta.env.PL_PATH_PREFIX}/`)
                  })
                  .catch(() => {
                    console.log('catched')
                    //form.setError('email', {message: 'invalid username and/or password'})
                    //form.setError('password', {message: 'invalid username and/or password'})

                    setShowWrongCredentialsTooltip(true)

                    setTimeout(() => {
                        setShowWrongCredentialsTooltip(false)
                      //form.clearErrors('email');
                      //form.clearErrors('password');
                    }, 2000)
                  })

                // 
              })}
              className="space-y-8"
            >
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input type="email" placeholder="E-Mail" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Password"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              {wrongUserOrPasswordTooltip(showWrongCredentialsTooltip)}
              <div className="flex w-full justify-center">
                <Button type="submit">Submit</Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
