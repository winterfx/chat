import { useEffect, useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { useNavigate } from 'react-router-dom';
import { User, Send, MessageSquare, LogOut, Trash2,ChevronUp } from 'lucide-react'
import { ScrollArea } from '@/components/ui/scroll-area';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import axios from 'axios';

import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarProvider ,SidebarMenu,SidebarMenuItem,SidebarMenuButton,SidebarTrigger} from '@/components/ui/sidebar'; // 假设你有 Sidebar 组件
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'; // 假设你有 DropdownMenu 组件

export default function Chat() {
  const navigate = useNavigate();
  interface Friend {
    uid: string;
    name: string;
    status: string;
  }

  interface ChatMessage {
    messageId: string;
    receiverId: string;
    senderId: string;
    data: string;
    timeStamp: number;
  }

  const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState<Record<string, ChatMessage[]>>({}); // 改为对象，以朋友的UID作为键
  const [friends, setFriends] = useState<Friend[]>([]); // 初始状态为空数组
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [userName, setUserName] = useState<string | null>(null);

  //the new machine's IP address is
  const host=import.meta.env.VITE_HOST
  const wsHost=import.meta.env.VITE_WS_HOST

  useEffect(() => {
    let intervalId;

    const initial = async () => {
      try {
        const response = await axios.get(`${host}/user/login-status`, {
          withCredentials: true // 确保请求包含凭据
        });
        const data = response.data as { uid: string | null; name: string | null };
        if (data.uid == null) {
          navigate('/login');
        } else {
          setUserId(data.uid);
          setUserName(data.name); // 存储用户名
          // 用户已登录，获取好友列表
          fetchFriends();
          connectWebSocket();
        }
        console.log('Login status:', data);
      } catch (error) {
        navigate('/login');
        console.error('Error checking login status:', error);
      }
    };

    const connectWebSocket = () => {
      const ws = new WebSocket(`${wsHost}/chat/ws`);
      ws.onopen = () => {
        console.log('WebSocket connection opened');
      };
      ws.onmessage = (event) => {
        const receivedMessage: ChatMessage = JSON.parse(event.data);
    
        // 确认消息的发送者和接收者
        const isForMe = receivedMessage.receiverId === userId; // 判断消息是否是发给我的
    
        // 更新消息状态
        setMessages((prevMessages) => {
            const friendMessages = prevMessages[receivedMessage.senderId] || []; // 获取发送者的消息列表
            return {
                ...prevMessages,
                [receivedMessage.senderId]: [...friendMessages, receivedMessage], // 更新发送者的消息
            };
        });
    
        console.log('WebSocket message received:', receivedMessage);
      };
      ws.onclose = () => {
        console.log('WebSocket connection closed');
      };
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      setWs(ws);
    };

    const fetchFriends = async () => {
      try {
        const response = await axios.get(`${host}/user/friends`, {
          withCredentials: true // 确保请求包含凭据
        });
        const data = response.data;
        setFriends(data); // 更新好友列表状态
        console.log('Friends data:', data);
      } catch (error) {
        console.error('Error fetching friends:', error);
      }
    };

    initial();
    intervalId = setInterval(fetchFriends, 15000); // 每15秒检查一次

    return () => clearInterval(intervalId); // 清除定时器
  }, [navigate]);

  const handleLogout = async () => {
    try {
      console.log('Logging out...');
      await axios.post(`${host}/user/logout`, {}, {
        withCredentials: true // 确保请求包含凭据
      });
      navigate('/login');
    } catch (error) {
      console.error('Error during logout:', error);
    }
  };

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!selectedFriend || !ws || !userId) {
      console.error('Invalid state for sending message');
      return;
    }

    const chatMessage: ChatMessage = {
      messageId: uuidv4(),
      receiverId: selectedFriend.uid,
      senderId: userId,
      data: message,
      timeStamp: Date.now(),
    };

    ws.send(JSON.stringify(chatMessage));
    setMessages((prevMessages) => {
      const friendMessages = prevMessages[selectedFriend.uid] || [];
      return {
          ...prevMessages,
          [selectedFriend.uid]: [...friendMessages, chatMessage],
      };
  });
    console.log(`Sending message to ${selectedFriend.name}: ${message}`);
    setMessage('');
  };

  return (
      <div className="flex h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
        <SidebarProvider>
          <SidebarTrigger />
            <Sidebar collapsible="icon" className="border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
              <SidebarHeader className="p-4 text-2xl font-bold bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">ChatApp</SidebarHeader>
              <SidebarContent className="px-2">
                  {friends.map((friend) => (
                      <button
                          key={friend.uid}
                          className={`flex items-center w-full p-3 my-1 rounded-lg transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-700 ${
                              selectedFriend?.uid === friend.uid ? 'bg-primary/10 dark:bg-primary/20' : ''
                          }`}
                          onClick={() => setSelectedFriend(friend)}
                      >
                          <div className="relative mr-3">
                              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center">
                                  <User className="h-5 w-5 text-primary" />
                              </div>
                              <span className={`absolute bottom-0 right-0 block h-3 w-3 rounded-full ring-2 ring-white dark:ring-gray-800 ${friend.status === "online" ? 'bg-green-400' : 'bg-gray-400'}`} />
                          </div>
                          <div className="flex flex-col items-start">
                              <span className="font-medium text-gray-800 dark:text-gray-100">{friend.name}</span>
                              <span className="text-xs text-gray-600 dark:text-gray-300">{friend.status}</span>
                          </div>
                      </button>
                  ))}
              </SidebarContent>
              <SidebarFooter className="border-t border-gray-200 dark:border-gray-700 p-4">
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <SidebarMenuButton className="w-full flex items-center p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center mr-2">
                              <User className="h-4 w-4 text-primary" />
                            </div>
                            <span className="font-medium text-gray-800 dark:text-gray-100">{userName}</span>
                            <ChevronUp className="ml-auto h-4 w-4 text-gray-600 dark:text-gray-300" />
                          </SidebarMenuButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-56">
                          <DropdownMenuItem onClick={handleLogout} className="text-red-600 dark:text-red-400">
                            <LogOut className="mr-2 h-4 w-4" />
                            <span>Log out</span>
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </SidebarMenuItem>
                  </SidebarMenu>
              </SidebarFooter>
            </Sidebar>
            <div className="flex-1 flex flex-col bg-white dark:bg-gray-900">
              {selectedFriend ? (
                  <>
                      <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md flex items-center">
                          <div className="flex items-center">
                              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center mr-3">
                                  <User className="h-5 w-5 text-primary" />
                              </div>
                              <div>
                                  <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200">{selectedFriend.name}</h2>
                                  <span className={`text-sm ${selectedFriend.status === "online" ? 'text-green-500' : 'text-gray-500'}`}>
                                      {selectedFriend.status}
                                  </span>
                              </div>
                          </div>
                      </div>
                      <ScrollArea className="flex-1 p-6">
                          <div className="space-y-4">
                              {messages[selectedFriend.uid]?.map((msg) => (
                                  <div key={msg.messageId} 
                                       className={`flex ${msg.senderId === userId ? 'justify-end' : 'justify-start'}`}>
                                      <div className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                                          msg.senderId === userId 
                                          ? 'bg-primary text-white' 
                                          : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-200'
                                      }`}>
                                          <p className="text-sm">{msg.data}</p>
                                          <span className="text-xs opacity-70 mt-1 block">
                                              {new Date(msg.timeStamp).toLocaleTimeString()}
                                          </span>
                                      </div>
                                  </div>
                              )) || (
                                  <div className="flex items-center justify-center h-full">
                                      <div className="text-center text-gray-500 dark:text-gray-400">
                                          <MessageSquare className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                          <p>No messages yet with {selectedFriend.name}</p>
                                          <p className="text-sm mt-2">Send a message to start the conversation</p>
                                      </div>
                                  </div>
                              )}
                          </div>
                      </ScrollArea>
                      <form onSubmit={handleSendMessage} 
                            className="p-4 border-t border-gray-200 dark:border-gray-700 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md">
                          <div className="flex items-center space-x-2">
                              <Input
                                  type="text"
                                  placeholder="Type your message..."
                                  value={message}
                                  onChange={(e) => setMessage(e.target.value)}
                                  className="flex-1 bg-gray-100 dark:bg-gray-800 border-0 focus-visible:ring-1 focus-visible:ring-primary"
                              />
                              <Button 
                                  type="submit" 
                                  className="bg-primary hover:bg-primary/90 text-white"
                                  disabled={!message.trim()}>
                                  <Send className="h-4 w-4" />
                              </Button>
                          </div>
                      </form>
                  </>
              ) : (
                  <div className="flex-1 flex flex-col items-center justify-center p-8">
                      <div className="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mb-6">
                          <MessageSquare className="h-10 w-10 text-primary" />
                      </div>
                      <h2 className="text-2xl font-bold text-gray-700 dark:text-gray-200 mb-2">Welcome to ChatApp</h2>
                      <p className="text-gray-500 dark:text-gray-400 text-center max-w-md">
                          Select a conversation from the sidebar to start chatting
                      </p>
                  </div>
              )}
            </div>
        </SidebarProvider>
      </div>
  );
}
