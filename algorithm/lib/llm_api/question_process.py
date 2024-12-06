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
        prompt = '''针对用户提出的问题，我将根据以下逻辑来判断是否需要进行互联网搜索以增强大模型的知识。请直接给出“是”或“否”的结果。
逻辑如下：
1. 如果问题属于如下类别回答否，否则回答是：
- 文章写作
- 代码调试
- 逻辑推理
- 数学计算问题

用户问题: {question}
最终回答：
'''
        use_rag = self.ask_llm(prompt=prompt.format(question=question),temperature=0.3)
        if '是' in use_rag:
            return True
        if '不' in use_rag:
            return False
        else:
            return False