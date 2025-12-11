document.addEventListener('DOMContentLoaded', function() {
    const messagesContainer = document.getElementById('messages');
    const userInput = document.getElementById('user-input');
    const sendButton = document.getElementById('send-btn');

    // Function to add a message to the chat
    function addMessage(text, isUser) {
        const messageDiv = document.createElement('div');
        messageDiv.classList.add('message');
        messageDiv.classList.add(isUser ? 'user-message' : 'ai-message');
        messageDiv.textContent = text;
        messagesContainer.appendChild(messageDiv);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
        return messageDiv;
    }

    // Function to show typing indicator
    function showTypingIndicator() {
        const typingDiv = document.createElement('div');
        typingDiv.id = 'typing-indicator';
        typingDiv.classList.add('typing-indicator');
        typingDiv.textContent = 'AI正在思考中...';
        messagesContainer.appendChild(typingDiv);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
        return typingDiv;
    }

    // Function to hide typing indicator
    function hideTypingIndicator() {
        const typingDiv = document.getElementById('typing-indicator');
        if (typingDiv) {
            typingDiv.remove();
        }
    }

    // Function to simulate typewriter effect with random speed for more natural feel
    function typeWriter(element, text, minSpeed = 10, maxSpeed = 50) {
        let i = 0;
        element.textContent = '';
        return new Promise(resolve => {
            function typeChar() {
                if (i < text.length) {
                    element.textContent += text.charAt(i);
                    i++;
                    messagesContainer.scrollTop = messagesContainer.scrollHeight;
                    
                    // Random delay for more natural typing effect
                    const randomDelay = Math.floor(Math.random() * (maxSpeed - minSpeed + 1)) + minSpeed;
                    setTimeout(typeChar, randomDelay);
                } else {
                    resolve();
                }
            }
            
            typeChar();
        });
    }

    // Function to send message to backend
    async function sendMessage() {
        const message = userInput.value.trim();
        if (!message) return;

        // Add user message to chat
        addMessage(message, true);
        userInput.value = '';

        // Disable input while waiting for response
        userInput.disabled = true;
        sendButton.disabled = true;

        // Show typing indicator
        const typingIndicator = showTypingIndicator();

        try {
            // Send message to backend
            const response = await fetch('/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ message: message })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            
            // Hide typing indicator
            hideTypingIndicator();
            
            // Add AI response to chat with typewriter effect
            const aiMessageElement = addMessage('', false);
            if (data.response) {
                await typeWriter(aiMessageElement, data.response);
            } else {
                await typeWriter(aiMessageElement, data.message || "抱歉，我没有理解您的问题。");
            }
        } catch (error) {
            console.error('Error:', error);
            hideTypingIndicator();
            const errorMessageElement = addMessage("抱歉，出现了错误，请稍后再试。", false);
        } finally {
            // Re-enable input
            userInput.disabled = false;
            sendButton.disabled = false;
            userInput.focus();
        }
    }

    // Event listeners
    sendButton.addEventListener('click', sendMessage);

    userInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    // Add welcome message
    const welcomeMessageElement = addMessage("", false);
    typeWriter(welcomeMessageElement, "您好！我是AI助手，有什么我可以帮您的吗？");
});