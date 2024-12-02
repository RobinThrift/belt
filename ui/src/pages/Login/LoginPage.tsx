import * as Form from "@/components/Form"
import { useBaseURL } from "@/hooks/useBaseURL"
import { useCSRFToken } from "@/hooks/useCSRFToken"
import { Belt, Password, User } from "@phosphor-icons/react"
import React from "react"

export interface LoginPageProps {
    redirectURL: string
    validationErrors?: {
        form?: string
    }
}

export function LoginPage(props: LoginPageProps) {
    let csrfToken = useCSRFToken()
    let baseURL = useBaseURL()

    return (
        <div className="h-screen w-screen flex items-center justify-stretch md:justify-center bg-body bg-[radial-gradient(#d2d5da_1px,transparent_1px)] [background-size:16px_16px] lg:[mask-image:radial-gradient(ellipse_50%_50%_at_50%_50%,#000_60%,transparent_100%)]">
            <div className="flex items-center justify-center flex-col w-full md:w-[400px]">
                <div className="flex flex-col items-center justify-center">
                    <Belt
                        className="text-primary animate-in slide-in-from-bottom fade-in-50 duration-300"
                        size={40}
                        weight="duotone"
                    />
                    <h1 className="text-primary text-3xl text-center font-semibold -mt-2 animate-in slide-in-from-bottom fade-in-50 duration-300">
                        Belt
                    </h1>
                    <h3 className="text-xl mt-8 mb-3 font-semibold text-primary-extra-dark animate-in slide-in-from-bottom fade-in-50 duration-300">
                        Sign In
                    </h3>
                </div>
                <Form.Root
                    action={`${baseURL}/login`}
                    method="post"
                    className="w-full"
                >
                    <input
                        type="hidden"
                        id="belt.csrf.token"
                        name="belt.csrf.token"
                        defaultValue={csrfToken}
                    />

                    <input
                        type="hidden"
                        id="redirect_url"
                        name="redirect_url"
                        defaultValue={props.redirectURL}
                    />

                    <div className="flex flex-col px-3 py-2 w-full">
                        <Form.Field
                            name="username"
                            aria-label="username"
                            className="input-field relative animate-in slide-in-from-bottom fade-in-50 duration-500"
                        >
                            <span className="icon ml-1 z-20">
                                <User />
                            </span>
                            <Form.Control asChild>
                                <input
                                    type="text"
                                    name="username"
                                    id="username"
                                    className="input rounded-b-none border-b-transparent py-2 pr-2 pl-10 hover:border-b-subtle-dark focus:border-b-subtle-dark focus:z-10 relative"
                                    autoComplete="username"
                                    placeholder="Username"
                                    required={true}
                                />
                            </Form.Control>
                        </Form.Field>

                        <Form.Field
                            name="password"
                            aria-label="password"
                            className="input-field relative animate-in slide-in-from-bottom fade-in-50 duration-500"
                        >
                            <span className="icon ml-1">
                                <Password />
                            </span>
                            <Form.Control asChild>
                                <input
                                    type="password"
                                    name="password"
                                    id="password"
                                    className="input rounded-t-none py-2 pr-2 pl-10"
                                    autoComplete="password"
                                    required={true}
                                    placeholder="Password"
                                />
                            </Form.Control>
                        </Form.Field>

                        <Form.Submit asChild>
                            <button
                                className="btn primary mt-5 animate-in slide-in-from-bottom fade-in-50 duration-500"
                                type="submit"
                            >
                                Login
                            </button>
                        </Form.Submit>
                    </div>

                    {props.validationErrors?.form && (
                        <div className="mt-5 field-message animate-in slide-in-from-bottom fade-in-50 duration-300">
                            {props.validationErrors.form}
                        </div>
                    )}
                </Form.Root>
            </div>
        </div>
    )
}
