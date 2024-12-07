from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI


class QuestionProcess (LLMBaseAPI):
    def split_question(self,question):
        # 对复杂问题进行拆分
        prompt='''现给你一个主问题，请参考如下示例，当主问题明确询问多个问题时，将主问题拆解成多个子问题，否则输出原始问题。输出要严格按照如下规定：
1. 子问题的每一项使用一行
2. 子问题的每一项是一个关键短句，注意主语明确，不可使用代词来指代前文对象
3. 子问题一定是主问题明确中提到的问题

示例主问题：给出太阳的形状或大小，以及工作原理或发光机制
子问题：
太阳的形状
太阳的大小
太阳的工作原理
太阳的发光机制

示例问题：太阳的发光原理
子问题：
太阳的发光原理

示例问题：斯坦福回应抄袭清华系大模型
子问题：
斯坦福回应抄袭清华系大模型

示例问题：实用水杯推荐
子问题：
实用水杯推荐

主问题：{question}
子问题：
''' 
        quesiton = self.ask_llm(prompt=prompt.format(question=question),temperature=0.2)
        quesiton_list = quesiton.split('\n')
        return quesiton_list
    
    # 当用户的输入过短的时候触发重写
    def question_rewrite_by_context(self,question,history):
        prompt = '''为了提高问题搜索的准确性和完整性，我需要你根据当前的提问（query）和对话上文（history）来重写问题。
请确保重写的问题能够准确反映用户的需求，并包含所有必要的信息，以便进行有效的搜索。重写后的问题不要超过30个词。
以下是当前的提问和对话上下文：
对话上文（history）: {history}
当前提问（query）: {question}

请重写问题，使其更加完整和方便搜索。'''
        quesiton_after_rewrite = self.ask_llm(prompt=prompt.format(question=question,history=self.messages_to_str(history)),temperature=0.2)
        return quesiton_after_rewrite

    def judge_use_rag(self,question):
        prompt = '''判断用户的请求或者描述是否属于以下类别：文章写作、代码调试、逻辑推理、数学计算问题。如果输入请回答是，不属于请回答否，只用回答是或者否一个字。

用户请求: {question}
回答：
'''
        use_rag = self.ask_llm(prompt=prompt.format(question=question),temperature=0.3)
        print(use_rag)
        if '是' in use_rag:
            return False
        if '不' in use_rag:
            return True
        else:
            return True
        
    def generate_more_related_question(self,question,answer,nums=3):
        prompt = '''根据用户的问题和回答，生成与该问题相关的问题或者话题，要求如下：
1. 生成的问题或话题的相关答案或描述，没有在回答中出现。
2. 生成的问题或话题更多的是帮助增加用户相关事情的了解和学习。
3. 生成的问题或话题要是用户基于当前上文可能会接下来想问的。
4. 生成每个问题或话题在单独一行

示例如下
用户问题：一句话介绍一下奥卡姆剃刀

回答：奥卡姆剃刀是一种由 14 世纪英国哲学家奥卡姆的威廉提出的原理，即 “如无必要，勿增实体”，简单来说就是在多种解释或假设中，应该优先选择最简单的那个。例如，在解释行星运动时，日心说相对地心说在模型上更简洁，就体现了奥卡姆剃刀的思维理念，它在科学、哲学等诸多领域被广泛用于筛选和评估理论。

生成{number}个相关的问题或话题：
奥卡姆剃刀原理的提出背景是什么？
有哪些科学理论的发展体现了奥卡姆剃刀原理？
在实际生活中如何运用奥卡姆剃刀原理？

真实任务
用户问题: {question}

回答：{answer}

生成{number}个相关的问题或话题：
'''     
        more_question_or_topic = self.ask_llm(prompt=prompt.format(question=question,answer=answer,number=nums),temperature=0.3)
        more_question_or_topic_list= more_question_or_topic.split('\n')
        return more_question_or_topic_list