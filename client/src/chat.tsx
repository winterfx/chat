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
      <div className="flex h-screen bg-background">
        <SidebarProvider >
          <SidebarTrigger />
            <Sidebar collapsible="icon">
              <SidebarHeader>ChatApp</SidebarHeader>
              <SidebarContent >
                  {friends.map((friend) => (
                      <button
                          key={friend.uid}
                          className={`flex items-center w-full p-4 hover:bg-primary-foreground/10 transition-colors ${
                              selectedFriend?.uid === friend.uid ? 'bg-primary-foreground/20' : ''
                          }`}
                          onClick={() => setSelectedFriend(friend)}
                      >
                          <div className="relative mr-2">
                              <User className="h-5 w-5" />
                              <span className={`absolute bottom-0 right-0 block h-2 w-2 rounded-full ring-2 ring-primary ${friend.status === "online" ? 'bg-green-400' : 'bg-gray-400'}`} />
                          </div>
                          <span className="flex-grow text-left">{friend.name}</span>
                      </button>
                  ))}
              </SidebarContent>
              <SidebarFooter>
              {/* <span className="mr-2">{userName}</span><Button onClick={handleLogout}>Log out</Button> */}
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <SidebarMenuButton>
                            <User /> {userName}
                            <ChevronUp className="ml-auto" />
                          </SidebarMenuButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={handleLogout}>
                            <LogOut className="mr-2 h-4 w-4" />
                            <span>Log out</span>
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={handleLogout}>
                            <Trash2 className="mr-2 h-4 w-4" />
                            <span>Delete account</span>
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </SidebarMenuItem>
                </SidebarMenu>
              </SidebarFooter>
            </Sidebar>
            <div className="flex-1 flex flex-col bg-secondary">
              {selectedFriend ? (
                  <>
                      <div className="bg-secondary-foreground text-secondary p-4 border-b border-secondary-foreground/10 flex items-center">
                          <h2 className="text-xl font-semibold flex-grow">{selectedFriend.name}</h2>
                          <span className={`text-sm ${selectedFriend.status === "online" ? 'text-green-400' : 'text-secondary/60'}`}>
                              {selectedFriend.status === "online" ? 'online' : 'offline'}
                          </span>
                      </div>
                      <ScrollArea className="flex-1 p-4">
                          {messages[selectedFriend.uid]?.map((msg) => (
                              <div key={msg.messageId} className={`mb-2 ${msg.senderId === userId ? 'text-right' : 'text-left'}`}>
                                  <p className={`text-secondary-foreground/70 ${msg.senderId === userId ? 'font-bold' : ''}`}>
                                      {msg.data}
                                  </p>
                              </div>
                          )) || (
                              <p className="text-secondary-foreground/70 text-center mt-4">
                                  还没有与 {selectedFriend.name} 的消息。
                              </p>
                          )}
                      </ScrollArea>
                      <form onSubmit={handleSendMessage} className="p-4 bg-secondary-foreground border-t border-secondary-foreground/10 flex">
                          <Input
                              type="text"
                              placeholder="Type a message..."
                              value={message}
                              onChange={(e) => setMessage(e.target.value)}
                              className="flex-1 mr-2 bg-secondary text-secondary-foreground"
                          />
                          <Button type="submit" variant="secondary">
                              <Send className="h-4 w-4 mr-2" />
                              Send
                          </Button>
                      </form>
                  </>
              ) : (
                  <div className="flex-1 flex flex-col items-center justify-center bg-secondary p-8 text-center">
                      <div className="bg-secondary-foreground rounded-full p-8 mb-6">
                          <MessageSquare className="h-16 w-16 text-secondary" />
                      </div>
                      <h2 className="text-3xl font-bold text-secondary-foreground mb-4">Welcome to ChatApp</h2>
                      <p className="text-secondary-foreground/70 max-w-md">
                          Select a friend from the list to start a conversation. Connect, chat, and stay in touch with your friends!
                      </p>
                  </div>
              )}
            </div>
        </SidebarProvider>
      </div>
  );
}
