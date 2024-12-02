from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI


class QuestionSplit (LLMBaseAPI):
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
        prompt = '''针对用户提出的问题，我将根据以下评估步骤来判断是否需要进行互联网搜索以增强大模型的知识。请在评估结束后，直接给出“是”或“否”的结果。
评估步骤：
1. **问题的性质**：是否为事实性问题，需要特定数据或信息。
2. **问题的时效性**：是否涉及最新的数据或近期事件。
3. **问题的复杂性**：是否需要复杂的推理或多个知识点的综合。
4. **问题的专业性**：是否涉及特定领域的专业知识。
5. **问题的上下文依赖性**：是否需要特定的上下文信息才能回答。
6. **大模型的能力**：大模型是否没有足够的知识范围和能力来回答这个问题。

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