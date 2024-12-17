import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import process from 'process';
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from '@/components/ui/button'
import { Icons } from '@/components/ui/icons'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import './index.css';

const host=import.meta.env.VITE_HOST
const wsHost=import.meta.env.VITE_WS_HOST

export default function Login() {
  const navigate = useNavigate();
  const handleGoogleSignIn = () => {
    window.location.href = `${host}/auth/google`;
  };

  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [isLoading, setIsLoading] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);

  useEffect(() => {
    const checkLoginStatus = async () => {
      try {
        const response = await axios.get(`${host}/user/login-status`, {
          withCredentials: true // 确保请求包含凭据
        });
        const data = response.data as { uid: string | null };
        if (data.uid != null) {
          navigate('/home'); // 如果已登录，重定向到主页
        }
        console.log('Login status:', data);
      } catch (error) {
        console.error('Error checking login status:', error);
      }
    };

    checkLoginStatus();
  }, []); // 仅在组件首次渲染时执行

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setIsLoading(true);
    try {
      const url = isSignUp ? `${host}/user/register` : `${host}/user/login`;
      const data = {
        email: email,
        password:password,
      };

      const response = await axios.post(url, data,{
        withCredentials: true // 确保请求包含凭据
      });
      console.log('Response:', response.data);
      navigate('/home');
    
    } catch (error) {
      console.error('Error:', error);
      // 处理错误响应，例如显示错误消息
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="relative flex items-center justify-center min-h-screen bg-gray-50 p-4 overflow-hidden">
      <div className="absolute inset-0 w-full h-full">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-600 opacity-[0.15]" />
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(90deg, transparent 0%, transparent 49%, rgba(0,0,0,0.05) 50%, transparent 51%, transparent 100%),
                           linear-gradient(0deg, transparent 0%, transparent 49%, rgba(0,0,0,0.05) 50%, transparent 51%, transparent 100%)`,
          backgroundSize: '30px 30px',
          animation: 'moveBg 30s linear infinite'
        }} />
      </div>
      <style jsx>{`
        @keyframes moveBg {
          0% {
            background-position: 0 0;
          }
          100% {
            background-position: 30px 30px;
          }
        }
      `}</style>
      <div className="w-full max-w-md transform transition-all duration-500 hover:scale-[1.02] relative z-10">
        <Card className="border-0 shadow-2xl bg-white/90 backdrop-blur-md">
          <CardHeader className="text-center space-y-4 pb-2">
            <div className="mx-auto mb-2 relative font-mono text-4xl text-primary select-none">
              <pre className="leading-[1.2]">
{`  ∩――∩
( ´•ω• )
/　　 づ
`}
              </pre>
            </div>
            <CardTitle className="text-4xl font-bold bg-gradient-to-r from-primary via-purple-500 to-pink-500 bg-clip-text text-transparent">
              Welcome Back
            </CardTitle>
            <CardDescription className="text-lg font-medium text-gray-600">Start your chat journey</CardDescription>
          </CardHeader>
          <CardContent className="pb-3">
            <div className="flex flex-col space-y-6">
              <Button 
                variant="outline" 
                className="w-full group flex items-center justify-center space-x-2 h-12 border-2 bg-white/50 backdrop-blur-sm hover:bg-primary/5 hover:border-primary/50 transition-all duration-300"
                onClick={handleGoogleSignIn}
                disabled={isLoading}
              >
                <Icons.react className="w-5 h-5 group-hover:scale-110 transition-transform duration-300" />
                <span>Continue with Google</span>
              </Button>
              
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t border-gray-300" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-white/80 backdrop-blur px-2 text-gray-500 font-medium">
                    Or continue with
                  </span>
                </div>
              </div>

              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-sm font-medium">Email</Label>
                  <Input
                    id="email"
                    placeholder="name@example.com"
                    type="email"
                    autoCapitalize="none"
                    autoComplete="email"
                    autoCorrect="off"
                    disabled={isLoading}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password" className="text-sm font-medium">Password</Label>
                  <Input
                    id="password"
                    placeholder="Enter your password"
                    type="password"
                    autoCapitalize="none"
                    autoComplete={isSignUp ? "new-password" : "current-password"}
                    disabled={isLoading}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </div>
                <Button 
                  type="submit" 
                  className="w-full h-12 text-lg font-semibold bg-gradient-to-r from-primary to-purple-600 hover:opacity-90 transition-all duration-300"
                  disabled={isLoading}
                >
                  {isLoading && <Icons.spinner className="mr-2 h-5 w-5 animate-spin" />}
                  {isSignUp ? "Create Account" : "Sign In"}
                </Button>
              </form>
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-4 text-center">
            <Button
              variant="ghost"
              className="text-sm text-gray-600 hover:text-primary transition-colors duration-300"
              onClick={() => setIsSignUp(!isSignUp)}
              disabled={isLoading}
            >
              {isSignUp ? "Already have an account? Sign in" : "Don't have an account? Sign up"}
            </Button>
            <p className="text-xs text-gray-500">
              By continuing, you agree to our{" "}
              <a href="/terms" className="text-primary hover:underline">Terms</a>
              {" "}and{" "}
              <a href="/privacy" className="text-primary hover:underline">Privacy Policy</a>
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}
